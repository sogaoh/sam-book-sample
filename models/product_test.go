package models

import (
	"sam-book-sample/db"

	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/guregu/dynamo"
)

func TestProduct_Create(t *testing.T) {
	db.SetupDBForTest()
	defer db.DropTable()

	// 新規作成用のデータを定義
	product := &Product{
		Name:        "テスト製品",
		Price:       100,
		ReleaseDate: time.Now(),
	}

	// Create メソッド呼び出し
	err := product.Create()
	assert.NoError(t, err)

	table, err := db.Table()
	assert.NoError(t, err)

	// 作成されたレコードを取得できるか、内容は合っているか確認
	var result ProductDynamo
	err = table.
		Get("PK", "Product-00000000001").
		Range("SK", dynamo.Equal, "00000000001").
		One(&result)
	assert.NoError(t, err)

	assert.Equal(t, product.Name, result.Name)
	assert.Equal(t, product.Price, result.Price)
	assert.Equal(t, product.ReleaseDate.Format("2006-01-02T15:04:05Z"), result.ReleaseDate.Format("2006-01-02T15:04:05Z"))
}


func TestProduct_Update(t *testing.T) {
	db.SetupDBForTest()
	defer db.DropTable()

	table, err := db.Table()
	assert.NoError(t, err)

	// 更新用レコードを作成
	createProductForTest(t, 1, "テスト製品", 100, time.Now())

	// 更新処理
	updateProduct := Product{
		Name:        "テスト製品（更新）",
		Price:       200,
		ReleaseDate: time.Now().AddDate(0, 0, 1),
		BaseModel:   BaseModel {
			ID:      1,
			Version: 1,
		},
	}
	err = updateProduct.Update(1)
	assert.NoError(t, err)

	//更新結果をチェック
	var result ProductDynamo
	err = table.
		Get("PK", "Product-00000000001").
		Range("SK", dynamo.Equal, "00000000001").
		One(&result)
	assert.NoError(t, err)

	assert.Equal(t, updateProduct.Name, result.Name)
	assert.Equal(t, updateProduct.Price, result.Price)
	assert.Equal(t, updateProduct.ReleaseDate.Format("2006-01-02T15:04:05Z"), result.ReleaseDate.Format("2006-01-02T15:04:05Z"))
}


// 取得用レコードを作成する関数
func createProductForTest(t *testing.T, id uint64, name string, price int, releaseDate time.Time) *ProductDynamo {
	t.Helper()

	table, err := db.Table()
	assert.NoError(t, err)

	product := ProductDynamo{
		MainTable: db.MainTable{
			PK: fmt.Sprintf("Product-%011d", id),
			SK: fmt.Sprintf("%011d", id),
		},
		Product: Product{
			Name:        name,
			Price:       price,
			ReleaseDate: releaseDate,
			BaseModel:   BaseModel{
				ID:      id,
				Version: 1,
			},
		},
	}

	err = table.Put(&product).Run()
	assert.NoError(t, err)

	return &product
}


func TestProduct_GetProducts(t *testing.T) {
	db.SetupDBForTest()
	defer db.DropTable()

	// 取得用レコードを作成
	product1 := createProductForTest(t, 1, "テスト製品1", 100, time.Now())
	product2 := createProductForTest(t, 2, "テスト製品2", 200, time.Now())

	products, err := GetProducts()
	assert.NoError(t, err)

	assert.Equal(t, product2.ID, products[0].ID)
	assert.Equal(t, product2.Name, products[0].Name)
	assert.Equal(t, product2.Price, products[0].Price)
	assert.Equal(t, product2.ReleaseDate.Format("2006-01-02T15:04:05Z"), products[0].ReleaseDate.Format("2006-01-02T15:04:05Z"))

	assert.Equal(t, product1.ID, products[1].ID)
	assert.Equal(t, product1.Name, products[1].Name)
	assert.Equal(t, product1.Price, products[1].Price)
	assert.Equal(t, product1.ReleaseDate.Format("2006-01-02T15:04:05Z"), products[1].ReleaseDate.Format("2006-01-02T15:04:05Z"))
}


func TestProduct_GetProductByID(t *testing.T) {
	db.SetupDBForTest()
	defer db.DropTable()

	// 取得用レコードを作成
	expected := createProductForTest(t, 1, "テスト製品1", 100, time.Now())

	product, err := GetProductByID(1)
	assert.NoError(t, err)

	assert.Equal(t, expected.ID, product.ID)
	assert.Equal(t, expected.Name, product.Name)
	assert.Equal(t, expected.Price, product.Price)
	assert.Equal(t, expected.ReleaseDate.Format("2006-01-02T15:04:05Z"), product.ReleaseDate.Format("2006-01-02T15:04:05Z"))
}


func TestProduct_DeleteProduct(t *testing.T) {
	db.SetupDBForTest()
	defer db.DropTable()

	table, err := db.Table()
	assert.NoError(t, err)

	// 削除用レコードを作成
	expected := createProductForTest(t, 1, "テスト製品", 100, time.Now())

	err = DeleteProduct(expected.ID)
	assert.NoError(t, err)

	// 削除されているかをチェック
	var result ProductDynamo
	err = table.
		Get("PK", "Product-00000000001").
		Range("SK", dynamo.Equal, "00000000001").
		One(&result)
	assert.Equal(t, dynamo.ErrNotFound, err)
}
