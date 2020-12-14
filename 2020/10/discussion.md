# First Test Input

8 possible jolt adapter combinations with this array of jolt adapters. (Array index on the top for reference.)

    [0 1 2 3 4 5  6  7  8  9 10 11 12]  <-- array index
    [0 1 4 5 6 7 10 11 12 15 16 19 22]  <-- jolt adapter array

The number of types of jolt differences across the array of all jolt adapters is:
    7 one-jolt drops
    5 three-jolt drops

7-element array of the index values of one-jolt drops

    [0 2 3 4 6 7 9]

array of substitutions

    [
        [2,3],
        [2,3,4],
        [3,4]
        [6,7]
    ]

A one-jolt drop can sometimes be substituted for a two-jolt drop or a three-jolt drop.
Take these one at a time. First look at a one-jolt drop > two-jolt drop transition.

one-jolt drop substituted by a two-jolt drop
----------------------------------------------
where there are two one-jolt drops consecutively, drop them both and the lowest index goes to the two-drop array

substitution 2A

    [0 1 2 3 4 5  6  7  8  9 10 11 12]
    [0 1 4   6 7 10 11 12 15 16 19 22]
        > one-drop: [2, 3]  > []
        > two-drop: []      > [2]
substitution 2B

    [0 1 2 3 4 5  6  7  8  9 10 11 12]
    [0 1 4 5   7 10 11 12 15 16 19 22]
        > one-drop: [3, 4]  > []
        > two-drop: []      > [3]
substitution 2C

    [0 1 2 3 4 5  6  7  8  9 10 11 12]
    [0 1 4 5 6 7 10    12 15 16 19 22]
        > one-drop: [6, 7]  > []
        > two-drop: []      > [6]

one-jolt drop substituted by a three-jolt drop
----------------------------------------------
where there are three one-drop jolts consecutively, drop them all and the lowest index goes to the three-drop array

substitution 3A

    [0 1 2 3 4 5  6  7  8  9 10 11 12]
    [0 1 4     7 10 11 12 15 16 19 22]
        > one-drop: [2, 3, 4]   > []
        > three-drop: []        > [2]
        > 2, 3, 4 > null

2A, 2B, and 3A are mutually exclusive
that is, only one of these substitutions can happen in any combination of jolt adapters:

    * 2A | one: [2, 3] > [] && two: [] > [2]
    * 2B | one: [3, 4] > [] && two: [] > [3]
    * 3A | one: [2, 3, 4] > [] && three: [] > [2]

We can determine that mathematically because the numbers in the "one-drop" arrays overlap.

    * [2, 3] overlaps with [3, 4] and [2, 3, 4]
    * [3, 4] overlaps with [2, 3] and [2, 3, 4]
    * [2, 3, 4] overlaps with [2, 3] and [3, 4]

If we didn't have substitution 2C, there would be only four possible combinations:

    1. the original
    2. 2A
    3. 2B
    4. 3A

But with substitution 2C, we have eight possible combinations:
- each of the four combinations above _without_ the 2C substitution
- each of the four combinations above _with_ the 2C substitution

number of possible combinations of substitutions
----------------------------------------------
    <index>      0   1  2  3  4  5   6   7   8   9  10  11   12
    1. original (0), 1, 4, 5, 6, 7, 10, 11, 12, 15, 16, 19, (22)
    2. 2C       (0), 1, 4, 5, 6, 7, 10,     12, 15, 16, 19, (22)
    3. 2B       (0), 1, 4, 5,    7, 10, 11, 12, 15, 16, 19, (22)
    4. 2B, 2C   (0), 1, 4, 5,    7, 10,     12, 15, 16, 19, (22)
    5. 2A       (0), 1, 4,    6, 7, 10, 11, 12, 15, 16, 19, (22)
    6. 2A, 2C   (0), 1, 4,    6, 7, 10,     12, 15, 16, 19, (22)
    7. 3A       (0), 1, 4,       7, 10, 11, 12, 15, 16, 19, (22)
    8. 3A, 2C   (0), 1, 4,       7, 10,     12, 15, 16, 19, (22)

In terms of jolt adapter substitutions (as opposed to the array of jolt adapters themselves, as shown above)

    1. original map[1:[0 2 3 4 6 7 9] 2:[]     3:[1 5 8 10 11]]
    2. 2C       map[1:[0 2 3 4 9]     2:[6]    3:[1 5 8 10 11]]
    3. 2B       map[1:[0 2 6 7 9]     2:[3]    3:[1 5 8 10 11]]
    4. 2B, 2C   map[1:[0 2 9]         2:[3 6]  3:[1 5 8 10 11]]
    5. 2A       map[1:[0 4 6 7 9]     2:[2]    3:[1 5 8 10 11]]
    6. 2A, 2C   map[1:[0 4 9]         2:[2 6]  3:[1 5 8 10 11]]
    7. 3A       map[1:[0 6 7 9]       2:[]     3:[1 2 5 8 10 11]]
    8. 3A, 2C   map[1:[0 7 9]         2:[6]    3:[1 2 5 8 10 11]]

To determine all this, 
1. we get an array of substitutions

        [
            [2, 3],
            [2, 3, 4],
            [3, 4],
            [6, 7]
        ]

2. we group the substitutions together into mutually exclusive clusters

        [
            # cluster with four possible combinations 
            # ([2, 3], [2, 3, 4], [3, 4], or None)
            [
                [2, 3],
                [2, 3, 4],
                [3, 4]
            ],
            
            # cluster with two possible combinations 
            # ([6, 7] or None)
            [6, 7]
        ]

3. compute the number of possible jolt adapter configurations.
    clusters with one substitution element have two possible combinations
    clusters with three substitution elements have four possible combinations
    clusters with five substitution elements have seven possible combinations
    
        2-combo clusters: 1 > 2^1 = 2
        4-combo clusters: 1 > 4^1 = 4
        7-combo clusters: 0 > 7^0 = 1

    multiply together to get the answer
    
        2 * 4 * 1 = 8 possible jolt adapter combinations

# Second Test Input

sorted jolt adapters from testinput2.txt

    [0 1 2 3 4 5 6 7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32]  <-- array index
    [0 1 2 3 4 7 8 9 10 11 14 17 18 19 20 23 24 25 28 31 32 33 34 35 38 39 42 45 46 47 48 49 52]  <-- jolt adapter array

22-element array of one-jolt drops (really 20-element array bc start at > 1

    [0 1 2 3 5 6 7 8 11 12 13 15 16 19 20 21 22 24 27 28 29 30]

First get an array of "substitutions" -- arrays of 2 or 3 consecutive integers

    [
        [0, 1],
        [0, 1, 2],
        [1, 2],
        [1, 2, 3],
        [2, 3],

        [5, 6],
        [5, 6, 7],
        [6, 7],
        [6, 7, 8],
        [7, 8],
        
        [11, 12],
        [11, 12, 13],
        [12, 13],
              
        [15, 16],      
        
        [19, 20],
        [19, 20, 21],
        [20, 21],
        [20, 21, 22],
        [21, 22],

        [27, 28]
        [27, 28, 29],
        [28, 29],
        [28, 29, 30],
        [29, 30]    
    ]
    
array of clusters of possible substitutions

    [
        # seven possible combos in this cluster
        [
            [0, 1],
            [0, 1, 2],
            [1, 2],
            [1, 2, 3],
            [2, 3],
        ],
        
        # seven possible combos in this cluster
        [
            [5, 6],
            [5, 6, 7],
            [6, 7],
            [6, 7, 8],
            [7, 8],
        ],
        
        # four possible combos in this cluster
        [
            [11, 12],
            [11, 12, 13],
            [12, 13]
        ],
        
        # two possible combos (either [15, 16] or none)
        [15, 16],
        
        # seven possible combos in this cluster
        [
            [19, 20],
            [19, 20, 21],
            [20, 21],
            [20, 21, 22],
            [21, 22]
        ],
        
        # seven possible combos in this cluster
        [
            [27, 28]
            [27, 28, 29],
            [28, 29],
            [28, 29, 30],
            [29, 30]
        ]
    ]

6 clusters

    1-combo clusters: 1 > 2^1  = 2      (clusters with len 1)
    4-combo clusters: 1 > 4^1  = 4      (clusters with len 3)
    7-combo clusters: 4 > 7^4  = 2401   (clusters with len 5)

multiply together to get the answer

    2 * 4 * 2401 = 19,208 possible jolt adapter combinations
