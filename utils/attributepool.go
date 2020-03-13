package utils

/*
  An AttributePool maintains a mapping from [key,value] Pairs called
  Attributes to Numbers (unsigened integers) and vice versa. These numbers are
  used to reference Attributes in Changesets.
*/

type AttributePool struct {
	NumToAttrib interface{} // e.g. {0: ['foo','bar']}
	AttribToNum interface{} // e.g. {'foo,bar': 0}
	NextNum     int
}

func (a *AttributePool) PutAttrib() {

}

func (a *AttributePool) GetAttrib() {

}

func (a *AttributePool) GetAttribKey() {

}

func (a *AttributePool) GetAttribValue() {

}

func (a *AttributePool) EachAttrib() {

}

func (a *AttributePool) ToJsonAble() {

}

func (a *AttributePool) FromJsonAble() {

}
