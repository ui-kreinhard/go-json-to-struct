What's this?
======
The idea is to generate the rough structs for marshalling and unmarshalling in go. A lot of time dealing with a rest api I have a sample json and have to model it with structs. This is teadious and error prone.

Now you can generate the rough structure with a simple cli. 

Some Notes: It takes the stdin and then outputs the result. It cannot handle arrays as top-level structure and datetime objects are generated to string objects.

What it currently can do
===
It can generate a rough structs from a json string including the annoying json tags.
Currently it can recognize sub structures, arrays, basic types. But dateitme is currently not supported

How to build/Use
===
Just run:
```
go build
```

Then run:
````
./go-json-to-struct
<< Pasete Json to stdin >>
````