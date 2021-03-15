#! /bin/bash
#Command examples
export RPC_ADDRESS=localhost:50051
echo $RPC_ADDRESS
numan #shows usage/help 
numan summary 
numan add 353-01-1234567  test.com "test carrier" #create number 
numan view 353-01-1234567 #check it 
numan allocate 353-01-1234567 55 #allocate to user 55
numan view 353-01-1234567 #check it 
numan list_user 55 #list by user 
numan list 353-01-1234567 #list by number 
numan list 353-01-12345 #list by partial number 
numan portout 353-01-1234567 30/9/2021 #set a porting out date
numan view 353-01-1234567 #check it 
numan portin 353-01-1234567 30/1/2021 #set a porting in date
numan deallocate 353-01-1234567 #de-allocate the number 

