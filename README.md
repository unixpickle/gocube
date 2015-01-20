# gocube

**gocube** will be a pure Go library for solving the 3x3x3 Rubik's cube and various parts of it.

# Design considerations

I don't know the best way to solve the cube on a computer. Thus, I will make this project as modular as possible so that it can be improved easily (unlike [some people's code](https://github.com/lgarron/shuang-chen-projects/blob/dee5de0485d20b6f7759e11b5aef248d9e0f2dda/min2phase-java/src/CoordCube.java#L225))

# TODO

The following things must be done:

 * Change search to send solutions on a channel
   * Allow more than one solution to be obtained
 * Implement phase-one solver
   * Implement corner goal
   * Create corner heuristic
 * Implement phase-two solver
   * Create heuristic 
