	
  == Notifier service  
	Use to send emails and text messages (email to phone number)  
	Uses a yaml configuration file to specify:  
	* smtp information  
	* tls cert and key files  
	* port to listen on  
	Takes a JSON payload with the format:  
	{  
		"Message" : "Message to Send",  
		"Recipients" : [  
			"user@email.com"  
		]  
	}  
