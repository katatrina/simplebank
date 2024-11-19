package util

const (
	DepositorRole = "depositor" // Normal users who deposit their money in the bank.
	BankerRole    = "banker"    // An employee of the bank who's in charge of customer service.
	/*
		 A depositor can only update their own information,
		while a banker can update the information of any users.
	*/
)
