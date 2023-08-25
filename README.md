Create a simple program in go to record the network traffic on the specified ports.

The program should receive the following arguments:

    1. Ports to listen
    2. Buffer size (limit for the structure in memory that store the network information)

Then, for each port, a goroutine must be created to catch the traffic of that port.

Each go routine should store the information in the shared map with the following definition:

`map[string][]string`

This map has the network port as the key and then we store the network traffic in the respective slice stored in that key.

The network trafic info must be a string with the following format:

"From 127.0.0.1 to 80 at 2023-08-25 10:23:02"

Note that the slices in the map must have the limit defined by the buffer size given by the program argument. When it reaches the limit, it should start saving from the beggining of the slice, overriding that way the oldest records.
