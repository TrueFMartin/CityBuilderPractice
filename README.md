# CityBuilderPractice
EngoEngine has a city builder tutorial. It is incomplete and hasn't been updated since Feb 2021. Trying to finish the incomplete(and most difficult) parts. 

Begining can be found here https://github.com/EngoEngine/TrafficManager/tree/01-hello-world

I have changed multiple parts for better code reusability. I have added an updated UI, new money system, created roads to be built between cities, new message systems, changed mouse functions, and more.

Current work:

I most recently added a path seeking system for roads built. Roads connecting multiple towns will upgrade to cities, to metros, ect.

Future work:

Next I will be adding audio, traffic managment, and traffic visuals. 

To run game, clone the repository, and in terminal run:

$go mod tidy

$go mod run .

This image shows a path between two cities. A path connecting to a city will recieve a teal highlight. A path that breaks off from that original path will have a red/oragne outline. Roads can be built with F1 and F2 at mouse position. Will later add curved roads and different colors for incomplete paths, wider paths, ect.

![CityBuilderExample](https://user-images.githubusercontent.com/103139765/210197409-ede4c54c-5b6e-4974-a095-9e72ff0424c0.png)

