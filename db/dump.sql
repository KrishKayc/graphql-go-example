USE BookCab;

Insert into Region (name, minLatitude, maxLatitude, minLongitude, maxLongitude) 
values("Adyar", 10.0000, 20.0000, 10.0000, 20.0000);

Insert into Region (name, minLatitude, maxLatitude, minLongitude, maxLongitude) 
values("Triplicane", 20.0000, 30.0000, 20.0000, 30.0000);

Insert into Location (name, latitude, longitude)
values("Indra Nagar", 12.6567, 15.3456);

Insert into Location (name, latitude, longitude)
values("Marina Beach", 23.5457, 25.3456);

Insert into Location (name, latitude, longitude)
values("Express Avenue", 21.5457, 26.3456);

-- Add user residing in Adyar
Insert into User(name, locationId, phoneNumber)
values("Raghu",1, 8220456032);

-- Add a driver residing in Marina Beach and triplicane
Insert into User(name, locationId, phoneNumber)
values("Krishna",2, 82205456033);

Insert into User(name, locationId, phoneNumber)
values("John",3, 82205456034);

-- Cab of the driver residing in marina beach, tenampet
Insert into Cab(type, driverId, number)
values (1, 2, "TN11AQ1234");

Insert into Cab(type, driverId, number)
values (1, 3, "TN11AQ1235");

