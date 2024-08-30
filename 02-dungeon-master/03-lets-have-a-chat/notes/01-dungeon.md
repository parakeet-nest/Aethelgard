# The Dungeon 🏰

I have a dungeon with 25 unique rooms in a 5*5 square map
Each room is a part of a larger square grid, with 5 rows and 5 columns:
Here's an ASCII representation of the dungeon with 5 rows and 5 columns, where each square is numbered room:

```
+---+---+---+---+---+
|  1|  2|  3|  4|  5|
+---+---+---+---+---+
|  6|  7|  8|  9| 10|
+---+---+---+---+---+
| 11| 12| 13| 14| 15|
+---+---+---+---+---+
| 16| 17| 18| 19| 20|
+---+---+---+---+---+
| 21| 22| 23| 24| 25|
+---+---+---+---+---+
```

Each room is labeled with its corresponding number, starting from 1 in the top-left corner and ending with 25 in the bottom-right corner.

Here is a list of possible connections between the rooms in the 5x5 grid, where each room can be connected to its adjacent rooms (up, down, left, right, and diagonally):

### Room Connections

- **Room 1:** Connected to Rooms 2, 6
- **Room 2:** Connected to Rooms 1, 3, 7
- **Room 3:** Connected to Rooms 2, 4, 8
- **Room 4:** Connected to Rooms 3, 5, 9
- **Room 5:** Connected to Rooms 4, 10
- **Room 6:** Connected to Rooms 1, 7, 11
- **Room 7:** Connected to Rooms 2, 6, 8, 12
- **Room 8:** Connected to Rooms 3, 7, 9, 13
- **Room 9:** Connected to Rooms 4, 8, 10, 14
- **Room 10:** Connected to Rooms 5, 9, 15
- **Room 11:** Connected to Rooms 6, 12, 16
- **Room 12:** Connected to Rooms 7, 11, 13, 17
- **Room 13:** Connected to Rooms 8, 12, 14, 18
- **Room 14:** Connected to Rooms 9, 13, 15, 19
- **Room 15:** Connected to Rooms 10, 14, 20
- **Room 16:** Connected to Rooms 11, 17, 21
- **Room 17:** Connected to Rooms 12, 16, 18, 22
- **Room 18:** Connected to Rooms 13, 17, 19, 23
- **Room 19:** Connected to Rooms 14, 18, 20, 24
- **Room 20:** Connected to Rooms 15, 19, 25
- **Room 21:** Connected to Rooms 16, 22
- **Room 22:** Connected to Rooms 17, 21, 23
- **Room 23:** Connected to Rooms 18, 22, 24
- **Room 24:** Connected to Rooms 19, 23, 25
- **Room 25:** Connected to Rooms 20, 24

This setup ensures that each room is connected to its immediate neighbors, allowing for movement throughout the entire grid.
