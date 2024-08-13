# Online store project

## In the current version, the server supports:


1) User registration. Request Example:  
    `host:port/user/sign-up`  
   `{
   "login": "username",
   "password": "userpassword"
   }`  
   The login must contain from 4 to 20 characters, the password from 6 to 20.

 
2) User authorization. Request Example:  
   `host:port/user/sign-in`  
   `{
   "login": "username",
   "password": "userpassword"
   }`


3) Changing the login or other user data. Request Example:  
    `host:port/user/update/:id`  
    `{
    "login": "username",
    "email": "email"
    }`  
   The login must contain from 4 to 20 characters.


4) Deleting the user. Request Example:  
   `host:port/user/delete/:id`
