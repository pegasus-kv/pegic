package pegic

/*
The AST flow. Drawed by asciiflow.com.

         |
         v
      SCAN
         +
         |
    +---------+
    v    |    v
DELETE   | COUNT
    +    |    +
    +---------+
         v
     HASHKEY:arg[1]
         +            +->SUFFIX:arg[1]+------+
         |            |                      |
         +-->SORTKEY--+->PREFIX:arg[1]+------+
         |            |                      |
         |            +->CONTAINS:arg[1]+----+
         |            |                      |
         |            +->BETWEEN:arg[1]      |
         |                  +                |
         |                  +-->AND:arg[1]+--+
         |                                   |
   +-----------------------------------------+
   v     v
 NOVALUE |
   +     |
   +-----+
         v
*/
