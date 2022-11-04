package pkg

/***********************************************************
 * Support class for serialization data parsing that holds
 * details of a single class to enable class data for that
 * class to be read (classDescFlags, field descriptions).
 *
 * Written by 51pwn
 **********************************************************/
type ClassDetails struct {
	_className         string        //The name of the class
	_refHandle         int           //The reference handle for the class
	_classDescFlags    uint8         //The classDescFlags value for the class
	_fieldDescriptions []*ClassField //The class field descriptions

}

func NewClassDetails(className string) *ClassDetails {
	return &ClassDetails{_className: className, _refHandle: -1, _classDescFlags: 0, _fieldDescriptions: []*ClassField{}}
}

/*******************
 * Get the class name.
 *
 * @return The class name.
 ******************/
func (this *ClassDetails) getClassName() string {
	return this._className
}

/*******************
 * Set the reference handle of the class.
 *
 * @param handle The reference handle value.
 ******************/
func (this *ClassDetails) setHandle(handle int) {
	this._refHandle = handle
}

/*******************
 * Get the reference handle.
 *
 * @return The reference handle value for this class.
 ******************/
func (this *ClassDetails) getHandle() int {
	return this._refHandle
}

/*******************
 * Set the classDescFlags property.
 *
 * @param classDescFlags The classDescFlags value.
 ******************/
func (this *ClassDetails) setClassDescFlags(classDescFlags uint8) {
	this._classDescFlags = classDescFlags
}

/*******************
 * Check whether the class is SC_SERIALIZABLE.
 *
 * @return True if the classDescFlags includes SC_SERIALIZABLE.
 ******************/
func (this *ClassDetails) isSC_SERIALIZABLE() bool {
	return (this._classDescFlags & 0x02) == 0x02
}

/*******************
 * Check whether the class is SC_EXTERNALIZABLE.
 *
 * @return True if the classDescFlags includes SC_EXTERNALIZABLE.
 ******************/
func (this *ClassDetails) isSC_EXTERNALIZABLE() bool {
	return (this._classDescFlags & 0x04) == 0x04
}

/*******************
 * Check whether the class is SC_WRITE_METHOD.
 *
 * @return True if the classDescFlags includes SC_WRITE_METHOD.
 ******************/
func (this *ClassDetails) isSC_WRITE_METHOD() bool {
	return (this._classDescFlags & 0x01) == 0x01
}

/*******************
 * Check whether the class is SC_BLOCK_DATA.
 *
 * @return True if the classDescFlags includes SC_BLOCK_DATA.
 ******************/
func (this *ClassDetails) isSC_BLOCK_DATA() bool {
	return (this._classDescFlags & 0x08) == 0x08
}

/*******************
 * Add a field description to the class details object.
 *
 * @param cf The ClassField object describing the field.
 ******************/
func (this *ClassDetails) addField(cf *ClassField) {
	this._fieldDescriptions = append(this._fieldDescriptions, cf)
}

/*******************
 * Get the class field descriptions.
 *
 * @return An array of field descriptions for the class.
 ******************/
func (this *ClassDetails) getFields() []*ClassField {
	return this._fieldDescriptions
}

/*******************
 * Set the name of the last field to be added to the ClassDetails object.
 *
 * @param name The field name.
 ******************/
func (this *ClassDetails) setLastFieldName(name string) {
	this._fieldDescriptions[len(this._fieldDescriptions)-1].setName(name)
}

/*******************
 * Set the className1 of the last field to be added to the ClassDetails
 * object.
 *
 * @param cn1 The className1 value.
 ******************/
func (this *ClassDetails) setLastFieldClassName1(cn1 string) {
	this._fieldDescriptions[len(this._fieldDescriptions)-1].setClassName1(cn1)
}
