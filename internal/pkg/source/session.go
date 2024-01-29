package source

type Session struct {
	classes map[JvmId]*JvmClass
	methods map[JvmId]*JvmMethod
	entries Entries
}

func (s *Session) Classes() map[JvmId]*JvmClass {
	return s.classes
}
func (s *Session) Methods() map[JvmId]*JvmMethod {
	return s.methods
}
func (s *Session) Entries() *Entries {
	return &s.entries
}

func (s *Session) GetArgType(id JvmId, index int) *string {
	if method := s.GetMethod(id); method != nil {
		return &method.a[index]
	}
	return nil
}

func (s *Session) GetReturnType(id JvmId) *string {
	if method := s.GetMethod(id); method != nil {
		return &method.r
	}
	return nil
}

func (s *Session) GetMethod(id JvmId) *JvmMethod {
	if s == nil {
		return nil
	}
	methods := s.Methods()
	if len(methods) == 0 {
		return nil
	}
	return methods[id]
}
