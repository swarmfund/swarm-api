package api

func (user *User) DocumentRejectReason() DocumentsRejectReasons {
	//entity := user.KYCEntities.GetSingle(KYCEntityTypeDocumentsRejectReasons)
	//var result DocumentsRejectReasons
	//result.Populate(entity)
	//return result
	panic("not implemented")
}

func (user *User) HaveDocument(entity int64, docType DocumentType) bool {
	if entity == 0 {
		return len(user.Documents[docType]) > 0
	}
	for _, doc := range user.Documents[docType] {
		if doc.EntityID == entity {
			return true
		}
	}
	return false
}

func (user *User) HaveEntityDoc(entity KYCEntity) bool {
	//map[doctype][]docs
	//requiredDocTypes := KYCRequiredDocs[entity.Type]
	//typeDocs := user.Documents[entity.Type]
	//for _, requiredDocType := range requiredDocTypes {
	//	ok := func(entityType KYCEntityType) bool {
	//		for _, doc := range user.Documents {
	//			doc
	//		}
	//	}(entity.Type)
	//}
	return false
}
