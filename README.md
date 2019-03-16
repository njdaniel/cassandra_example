# cassandra-example 

Setup Cassandra:

$ brew install ccm maven

$ pip install cqlsh

$ ccm create -v 3.9 streamdemoapi

$ ccm populate -n 1

$ ccm start

$ echo "CREATE KEYSPACE streamdemoapi WITH \ replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};" | cqlsh --cqlversion 3.4.2

$ echo " use streamdemoapi; create table messages ( id UUID, user_id UUID, Message text, PRIMARY KEY(id) );" | cqlsh --cqlversion 3.4.2

$ echo " use streamdemoapi; CREATE TABLE users ( id UUID, firstname text, lastname text, age int, email text, city text, PRIMARY KEY (id) );" | cqlsh --cqlversion 3.4.2

Testing Create new user: curl -X POST \ -H "Content-Type: application/x-www-form-urlencoded" \ -d 'firstname=Ian&lastname=Douglas&city=Boulder&email=ian@getstream.io&age=42' \ "http://localhost:8080/users/new"