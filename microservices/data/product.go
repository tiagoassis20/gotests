package domain

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
}

type Products []*Product

var listProduct = []*Product{
	&Product{0, "kitkat", "chocolate", 2.5},
	&Product{1, "bis", "chocolate", 3.5},
	&Product{2, "leite tirol", "leite integral", 2.0},
}

func GetProducts() (Products, error) {
	return listProduct, nil
}

func GetProduct(id int) (*Product, error) {
	return listProduct[id], nil
}

func AddProduct(p *Product) (*Product, error){
	p.id
	return p
}
