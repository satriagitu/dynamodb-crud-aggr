# DynamoDB with AWS SDK for Go

## 🔹 Persiapan

### 1️⃣ Install AWS SDK for Go
Jalankan perintah berikut untuk menginstall AWS SDK:

```sh
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/dynamodb
```

### 2️⃣ Jalankan DynamoDB Lokal
Untuk menjalankan DynamoDB lokal menggunakan Docker, gunakan perintah berikut:

```sh
docker run -p 8000:8000 amazon/dynamodb-local
```

---

## 🔹 Buat Tabel Orders dengan Global Secondary Index (GSI)
Gunakan perintah berikut untuk membuat tabel **Orders** dengan **Global Secondary Index (GSI)**:

```sh
aws dynamodb create-table \
    --table-name Orders \
    --attribute-definitions AttributeName=OrderID,AttributeType=S AttributeName=CustomerID,AttributeType=S \
    --key-schema AttributeName=OrderID,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --global-secondary-indexes '[
        {
            "IndexName": "CustomerIndex",
            "KeySchema": [{"AttributeName": "CustomerID", "KeyType": "HASH"}],
            "Projection": {"ProjectionType": "ALL"}
        }
    ]' \
    --endpoint-url http://localhost:8000
```

Dengan konfigurasi di atas:
- **`OrderID`** sebagai Primary Key (**HASH**).
- **`CustomerID`** sebagai atribut tambahan untuk **GSI**.
- **Billing mode** menggunakan `PAY_PER_REQUEST`.
- **GSI `CustomerIndex`** memungkinkan query berdasarkan `CustomerID`.
