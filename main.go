package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Order struct {
	OrderID    string
	CustomerID string
	Amount     int
	Status     string
}

var svc *dynamodb.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(aws.EndpointResolverFunc(
		func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "http://localhost:8000"}, nil
		},
	)))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	svc = dynamodb.NewFromConfig(cfg)
}

// 1️⃣ CREATE DATA (Batch Insert)
func batchInsertOrders() {
	orders := []Order{
		{"O1", "C1", 100, "Completed"},
		{"O2", "C1", 200, "Pending"},
		{"O3", "C2", 150, "Completed"},
		{"O4", "C3", 300, "Shipped"},
	}

	writeReqs := []types.WriteRequest{}
	for _, order := range orders {
		item := map[string]types.AttributeValue{
			"OrderID":    &types.AttributeValueMemberS{Value: order.OrderID},
			"CustomerID": &types.AttributeValueMemberS{Value: order.CustomerID},
			"Amount":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.Amount)},
			"Status":     &types.AttributeValueMemberS{Value: order.Status},
		}
		writeReqs = append(writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}

	_, err := svc.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{"Orders": writeReqs},
	})
	if err != nil {
		log.Fatalf("Batch insert failed: %v", err)
	}
	fmt.Println("✅ Orders inserted successfully!")
}

// 2️⃣ READ DATA (Get All Orders)
func getAllOrders() {
	resp, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Orders"),
	})
	if err != nil {
		log.Fatalf("Failed to get orders: %v", err)
	}

	fmt.Println("📌 All Orders:")
	for _, item := range resp.Items {
		fmt.Printf("OrderID: %s, CustomerID: %s, Amount: %s, Status: %s\n",
			item["OrderID"].(*types.AttributeValueMemberS).Value,
			item["CustomerID"].(*types.AttributeValueMemberS).Value,
			item["Amount"].(*types.AttributeValueMemberN).Value,
			item["Status"].(*types.AttributeValueMemberS).Value,
		)
	}
}

// 3️⃣ UPDATE DATA (Update Status)
func updateOrderStatus(orderID, newStatus string) {
	_, err := svc.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("Orders"),
		Key: map[string]types.AttributeValue{
			"OrderID": &types.AttributeValueMemberS{Value: orderID},
		},
		UpdateExpression: aws.String("SET #s = :newStatus"),
		ExpressionAttributeNames: map[string]string{
			"#s": "Status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newStatus": &types.AttributeValueMemberS{Value: newStatus},
		},
	})
	if err != nil {
		log.Fatalf("Failed to update order: %v", err)
	}
	fmt.Printf("✅ Order %s updated to %s\n", orderID, newStatus)
}

// 4️⃣ DELETE DATA (Delete Order)
func deleteOrder(orderID string) {
	_, err := svc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Orders"),
		Key: map[string]types.AttributeValue{
			"OrderID": &types.AttributeValueMemberS{Value: orderID},
		},
	})
	if err != nil {
		log.Fatalf("Failed to delete order: %v", err)
	}
	fmt.Printf("✅ Order %s deleted successfully!\n", orderID)
}

// 5️⃣ AGGREGATION (Total Amount by Customer)
func aggregateTotalAmountByCustomer(customerID string) {
	resp, err := svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String("Orders"),
		IndexName: aws.String("CustomerIndex"),
		KeyConditions: map[string]types.Condition{
			"CustomerID": {
				ComparisonOperator: types.ComparisonOperatorEq,
				AttributeValueList: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: customerID},
				},
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to aggregate: %v", err)
	}

	total := 0
	for _, item := range resp.Items {
		amount := item["Amount"].(*types.AttributeValueMemberN).Value
		var value int
		fmt.Sscanf(amount, "%d", &value)
		total += value
	}
	fmt.Printf("📊 Total amount spent by Customer %s: %d\n", customerID, total)
}

func main() {
	batchInsertOrders()
	getAllOrders()
	updateOrderStatus("O2", "Shipped")
	aggregateTotalAmountByCustomer("C1")

	deleteOrder("O4")
	getAllOrders()
}
