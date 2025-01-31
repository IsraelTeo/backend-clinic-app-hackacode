package auth

func comparePassword(loginPassword, userPassword string) bool {
	return loginPassword == userPassword
}
