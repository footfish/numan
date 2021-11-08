#! /bin/bash
#Command examples
export RPC_ADDRESS=localhost:50051
echo $RPC_ADDRESS
num #shows usage/help 
num summary 
num add 353-01-1231234  test.com "test carrier" #create number 
num list 353-01-1231234 #check it 
num portin 353-01-1231234 30/1/2021 #set a porting in date
num allocate 353-01-1231234 55 #allocate to ownerID 55
num owner 55 #list by owner
num list 353-01-1231234 #list by number 
num list 353-01-12345 #list by partial number 
num portout 353-01-1231234 30/9/2021 #set a porting out date
num deallocate 353-01-1231234 55 #de-allocate the number 
num list 353-01-1231234 #check it (currently port out only shown after deallocation) 

