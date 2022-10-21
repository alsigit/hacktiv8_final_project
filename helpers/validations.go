package helpers

func User_Age(age int) (bool, string) {
	if age == 0 {
		return false, "Age cannot be empty."
	} else if age <= 8 {
		return false, "Minimum age is 8."
	}
	return true, ""
}

func User_Email(email string) (bool, string) {
	if email == "" {
		return false, "Email cannot be empty."
	} else if !ValidMailAddress(email) {
		return false, "Please provide correct email address."
	}
	return true, ""
}

func User_Username(username string) (bool, string) {
	if username == "" {
		return false, "Username cannot be empty."
	}
	return true, ""
}

func User_Password(password string) (bool, string) {
	if password == "" {
		return false, "Password cannot be empty."
	} else if len(password) < 6 {
		return false, "Password must be 6 or more character(s)."
	}
	return true, ""
}

func Photos_Title(title string) (bool, string) {
	if title == "" {
		return false, "Title cannot be empty."
	}
	return true, ""
}

func Photos_Url(url string) (bool, string) {
	if url == "" {
		return false, "Photo URL cannot be empty."
	}
	return true, ""
}

func Comments_Message(msg string) (bool, string) {
	if msg == "" {
		return false, "Message cannot be empty."
	}
	return true, ""
}

func SocMed_Name(msg string) (bool, string) {
	if msg == "" {
		return false, "Social media name cannot be empty."
	}
	return true, ""
}

func SocMed_URL(msg string) (bool, string) {
	if msg == "" {
		return false, "Social media url cannot be empty."
	}
	return true, ""
}
