# CityBuilderPractice
EngoEngine has a city builder tutorial. It is incomplete and hasn't been updated since Feb 2021. Trying to finish the incomplete(and most difficult) parts. 

Beginning can be found here https://github.com/EngoEngine/TrafficManager/tree/01-hello-world

I have changed multiple parts for better code re-usability. I have added an updated UI, new money system, created roads to be built between cities, new message systems, changed mouse functions, and more.

## Current work:

I most recently added a path seeking system for roads built. Roads connecting multiple towns will upgrade to cities, to metros, ect.

### 1/8/2023: 

Currently every road built not at the front of a current path creates new paths. For example Image the following grid(R = road, C = city, NR = new road):

    _, 1, 2, 3, 4, 5, 6
    
    1, r, r, --, --, --, --
    
    2, --, r, nr, --, --,--
    
    3, r, r, --, --, --, --
    
    4, --, r, --, --, --, --
    
    before the new road at Row:3, Col: 2 there are two paths:
    
    [1][1] -> [1][2] -> [2][2] -> [3][2] -> [4][2] &&
    
    [1][1] -> [1][2] -> [2][2] -> [3][2] -> [3][1]
    
    After the new road at [3][2], there is obviously a new path from [1][1] -> -> -> [2][3]. 
    But is there a path from: 
    
    [3][1] -> -> -> [2][3]? 
    
    What about:   
    
    [4][2] -> -> -> [2][3]? 
    
All of these paths must be considered if a new city later appears adjacent to one of these roads. 

The **Solution**: All paths must begin at a city.

To solve this I am changing it so roads can only be built adjacent to a city,
or adjacent to an already present path. This will significantly reduce the number of paths that must be considered.



## Future work:

    Next I will be adding audio, traffic management, and traffic visuals. 

## How to use:
To run game, clone the repository, and in terminal run:

    $go mod tidy
    
    $go mod run .

This image shows a path between two cities. A path connecting to a city will receive a teal highlight.
A path that breaks off from that original path will have a red/orange outline. 
Roads can be built by clicking adjacent to a city or road while you have >= $50. Will later add different colors for incomplete paths, wider paths, ect.

![CityBuilderExample](https://user-images.githubusercontent.com/103139765/210197409-ede4c54c-5b6e-4974-a095-9e72ff0424c0.png)

