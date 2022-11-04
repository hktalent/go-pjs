package pkg

/***********************************************************
 * Support class for serialization data parsing that holds
 * details of a class field to enable the field value to
 * be read from the stream.
 *
 * Written by 51pwn
 **********************************************************/
type ClassField struct {
	_typeCode   uint8  `json:"_TypeCode"`   // The field type code
	_name       string `json:"_Name"`       // The field name
	_className1 string `json:"_ClassName1"` // The className1 property (object and array type fields)
}

/*******************
 * Construct the ClassField object.
 *
 * @param typeCode The field type code.
 ******************/
func NewClassField(typeCode uint8) *ClassField {
	return &ClassField{_typeCode: typeCode, _name: "", _className1: ""}
}

/*******************
 * Get the field type code.
 *
 * @return The field type code.
 ******************/
func (this *ClassField) getTypeCode() uint8 {
	return this._typeCode
}

/*******************
 * Set the field name.
 *
 * @param name The field name.
 ******************/
func (this *ClassField) setName(name string) {
	this._name = name
}

/*******************
 * Get the field name.
 *
 * @return The field name.
 ******************/
func (this *ClassField) getName() string {
	return this._name
}

/*******************
 * Set the className1 property of the field.
 *
 * @param cn1 The className1 value.
 ******************/
func (this *ClassField) setClassName1(cn1 string) {
	this._className1 = cn1
}
