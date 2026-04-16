package firestore

func (f *Firestore) Close() error {
	if f.FirestoreClient != nil {
		return f.FirestoreClient.Close()
	}
	return nil
}
