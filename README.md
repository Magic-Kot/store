# Online store project

## In the current version, the server supports:


1) User registration. Request Example:  
    `host:port/auth/sign-up`  
   `{
   "login": "username",
   "password": "userpassword"
   }`  
   The login must contain from 4 to 20 characters, the password from 6 to 20.

 
2) User authorization. Request Example:  
   `host:port/auth/sign-in`  
   `{
   "login": "username",
   "password": "userpassword"
   }`

   
3) Getting user data. Request Example:
   `host:port/user/get`  
   `Header. Authorization: token`


4) Changing the login or other user data. Request Example:  
    `host:port/user/update`  
    `{
    "login": "username",
    "email": "email"
    }`  
   `Header. Authorization: token`  
   The login must contain from 4 to 20 characters.


5) Deleting the user. Request Example:  
   `host:port/user/delete`
   `Header. Authorization: token` 
