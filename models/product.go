package models

import (
	"sam-book-sample/db"
	"github.com/memememomo/nomof"

	"time"

	"github.com/pkg/errors"
	"github.com/guregu/dynamo"
)

type Product struct {
	BaseModel
	Name        string    `dynamo:"Name"`
	Price       int       `dynamo:"Price"`
	ReleaseDate time.Time `dynamo:"ReleaseDate"`
}

type ProductDynamo struct {
	db.MainTable
	Product
}

// implements DynamoModelMapper

func (p *Product) EntityName() string {
	return getEntityNameFromStruct(*p)
}

func (p *Product) PK() string {
	return getPK(p)
}

func (p *Product) SK() string {
	return getSK(p)
}

func (p *Product) PutToDynamo() error {
	return putEntityToDynamo(p)
}

func (p *Product) CreateDynamoRecord() error {
	return createEntityToDynamo(p)
}

func (p *Product) UpdateDynamoRecord() error {
	return updateEntityToDynamo(p)
}

func (p *Product) DeleteDynamoRecord() error {
	return deleteEntity(p)
}

func (p *Product) SetID(id uint64) {
	p.ID = id
}

func (p *Product) GetID() uint64 {
	return p.ID
}

func (p *Product) SetVersion(v int) {
	p.Version = v
}

func (p *Product) GetVersion() int {
	return p.Version
}

func (p *Product) GenerateRecord() interface{} {
	return &ProductDynamo{
		MainTable: generateMainTable(p),
		Product:   *p,
	}
}


func (p *Product) Create() error {
	return p.PutToDynamo()
}

func (p *Product) Update(id uint64) error {
	p.ID = id
	return p.PutToDynamo()
}


func GetProducts() ([]*Product, error) {
	table, err := db.Table()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// PK カラムの値が”Product-”で始まるデータを取得するフィルタを作成
	fb := nomof.NewBuilder()
	fb.BeginsWith(db.PKName, (&Product{}).EntityName())

	// Scan でフィルタを渡して対象となるデータを取得する
	var productDynamo []ProductDynamo
	err = table.
			Scan().
			Filter(fb.JoinAnd(), fb.Arg...).
			All(&productDynamo)

	if err != nil {
		if err == dynamo.ErrNotFound {
			return []*Product{}, nil
		}
	}

	// []ProductDynamo から []*Product に変換する
	var products = make([]*Product, len(productDynamo))
	for i := 0; i < len(productDynamo); i++ {
		products[i] = & productDynamo[i].Product
	}

	return products, nil
}


func GetProductByID(id uint64) (*Product, error) {
	var product ProductDynamo
	ret, err := getEntityByID(id, &Product{}, &product)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if ret == nil {
		return nil, nil
	}
	return &product.Product, nil
}


func DeleteProduct(id uint64) error {
	product, err := GetProductByID(id)
	if err != nil {
		return errors.WithStack(err)
	}
	if product == nil {
		return nil
	}
	err = deleteEntity(product)

	return errors.WithStack(err)
}
