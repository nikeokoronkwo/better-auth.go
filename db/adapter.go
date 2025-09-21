package db

type AdapterFactoryCustomizeAdapterCreator func (config any) CustomAdapter

type CustomAdapter interface {

} 
	

type Adapter interface {

} 

type AdapterOptions struct {
	Config struct{

	}
	Adapter AdapterFactoryCustomizeAdapterCreator
}

func CreateAdapter(opts AdapterOptions) Adapter {

}