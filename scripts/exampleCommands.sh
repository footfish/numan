#! /bin/bash
#Command examples
export RPC_ADDRESS=localhost:50051
echo $RPC_ADDRESS
num #shows usage/help 
num summary 
num add 353-01-1234568  test.com "test carrier" #create number 
num view 353-01-1234568 #check it 
num portin 353-01-1234568 30/1/2021 #set a porting in date
num allocate 353-01-1234568 55 #allocate to ownerID 55
num view 353-01-1234568 #check it 
num list_owner 55 #list by owner
num list 353-01-1234568 #list by number 
num list 353-01-12345 #list by partial number 
num portout 353-01-1234568 30/9/2021 #set a porting out date
num view 353-01-1234568 #check it (currently port out only shown after deallocation) 
num deallocate 353-01-1234568 #de-allocate the number 
num view 353-01-1234568 #check it (currently port out only shown after deallocation) 

