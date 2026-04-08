package env

import "os"

func Init() *Env {
	return &Env{
		FirestoreProjectID:      os.Getenv("FIRESTORE_PROJECT_ID"),
		FirestoreUserCollection: os.Getenv("FIRESTORE_COLLECTION_USERS"),
	}
}
