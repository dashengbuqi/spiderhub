package widgets

type Combobox struct {
	Id       string
	Name     string
	Required bool
	Editable bool
	Data     string
	Value    interface{}
	Width    int
	Multiple bool
	OnChange string
}
