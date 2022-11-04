package pkg

/***********************************************************
 * Support class for serialization data parsing that holds
 * class data details that are required to enable object
 * properties and array elements to be read from the
 * stream.
 *
 * Written by 51pwn
 **********************************************************/
type ClassDataDesc struct {
	/*******************
	 * Properties
	 ******************/
	_classDetails []*ClassDetails // //List of all classes making up this class data description (i.e. class, super class, etc)
}

/*******************
 * Construct the class data description object.
 ******************/
func NewClassDataDesc() *ClassDataDesc {
	return &ClassDataDesc{_classDetails: []*ClassDetails{}}
}

func NewClassDataDesc1(cd []*ClassDetails) *ClassDataDesc {
	return &ClassDataDesc{_classDetails: cd}
}

/*******************
 * Private constructor which creates a new ClassDataDesc and initialises it
 * with a subset of the ClassDetails objects from another.
 *
 * @param cd The list of ClassDetails objects for the new ClassDataDesc.
 ******************/
func (this *ClassDataDesc) ClassDataDesc(cd []*ClassDetails) {
	this._classDetails = cd
}

func (this *ClassDataDesc) GetClassDataDesc() []*ClassDetails {
	return this._classDetails
}

/*******************
 * Build a new ClassDataDesc object from the given class index.
 *
 * This is used to enable classdata to be read from the stream in the case
 * where a classDesc element references a super class of another classDesc.
 *
 * @param index The index to start the new ClassDataDesc from.
 * @return A ClassDataDesc describing the classes from the given index.
 ******************/
func (this *ClassDataDesc) buildClassDataDescFromIndex(index int) *ClassDataDesc {
	var cd []*ClassDetails
	// Build a list of the ClassDetails objects for the new ClassDataDesc
	cd = append(cd, this._classDetails...)

	//Return a new ClassDataDesc describing this subset of classes
	return NewClassDataDesc1(cd)
}

/*******************
 * Add a super class data description to this ClassDataDesc by copying the
 * class details across to this one.
 *
 * @param scdd The ClassDataDesc object describing the super class.
 ******************/
func (this *ClassDataDesc) addSuperClassDesc(scdd *ClassDataDesc) {
	// Copy the ClassDetails elements to this ClassDataDesc object
	if scdd != nil {
		this._classDetails = append(this._classDetails, scdd.GetClassDataDesc()...)
	}
}

/*******************
 * Add a class to the ClassDataDesc by name.
 *
 * @param className The name of the class to add.
 ******************/
func (this *ClassDataDesc) addClass(className string) {
	this._classDetails = append(this._classDetails, NewClassDetails(className))
}

/*******************
 * Set the reference handle of the last class to be added to the
 * ClassDataDesc.
 *
 * @param handle The handle value.
 ******************/
func (this *ClassDataDesc) setLastClassHandle(handle int) {
	this._classDetails[len(this._classDetails)-1].setHandle(handle)
}

/*******************
 * Set the classDescFlags of the last class to be added to the
 * ClassDataDesc.
 *
 * @param classDescFlags The classDescFlags value.
 ******************/
func (this *ClassDataDesc) setLastClassDescFlags(classDescFlags uint8) {
	this._classDetails[len(this._classDetails)-1].setClassDescFlags(classDescFlags)
}

/*******************
 * Add a field with the given type code to the last class to be added to
 * the ClassDataDesc.
 *
 * @param typeCode The field type code.
 ******************/
func (this *ClassDataDesc) addFieldToLastClass(typeCode uint8) {
	this._classDetails[len(this._classDetails)-1].addField(NewClassField(typeCode))
}

/*******************
 * Set the name of the last field that was added to the last class to be
 * added to the ClassDataDesc.
 *
 * @param name The field name.
 ******************/
func (this *ClassDataDesc) setLastFieldName(name string) {
	this._classDetails[len(this._classDetails)-1].setLastFieldName(name)
}

/*******************
 * Set the className1 of the last field that was added to the last class to
 * be added to the ClassDataDesc.
 *
 * @param cn1 The className1 value.
 ******************/
func (this *ClassDataDesc) setLastFieldClassName1(cn1 string) {
	this._classDetails[len(this._classDetails)-1].setLastFieldClassName1(cn1)
}

/*******************
 * Get the details of a class by index.
 *
 * @param index The index of the class to retrieve details of.
 * @return The requested ClassDetails object.
 ******************/
func (this *ClassDataDesc) getClassDetails(index int) *ClassDetails {
	if index >= len(this._classDetails) {
		return nil
	}
	return this._classDetails[index]
}

/*******************
 * Get the number of classes making up this class data description.
 *
 * @return The number of classes making up this class data description.
 ******************/
func (this *ClassDataDesc) getClassCount() int {
	return len(this._classDetails)
}
