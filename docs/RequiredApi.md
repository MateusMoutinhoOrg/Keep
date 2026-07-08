
### Write (key: string, value: Byte[]) (error or null):
#### Description
Write Bytes to the key
#### Args:
  - key:The Key to insert
  - value:The Bytes to insert
#### Returns:
  - error:Error if happened something

----
### WriteIfKeyNotExists(key: string, value: Byte[]) (error or null):
#### Description
Write Bytes to the key if the key not exists
#### Args:
  - key:The Key to insert
  - value:The Bytes to insert
#### Returns:
  - error:Error if happened something

----
### WriteifValueEquals(key: string, value: Byte[], oldValue: Byte[]) (error or null):
#### Description
Write Bytes to the key if the value equals oldValue
#### Args:
  - key:The Key to insert
  - value:The Bytes to insert
  - oldValue:The Bytes to compare
#### Returns:
  - error:Error if happened something
----
### Append(key:string, value:Byte[]) (error or null):
#### Description
Append Bytes to the key
#### Args:
  - key:The Key to append
  - value:The Bytes to append
#### Returns:
  - error:Error if happened something
----

### InsertAt(key:string, position:int64, value:Byte[]) (error or null):
#### Description
Insert Bytes to the key at the specified position
#### Args:
  - key:The Key to insert
  - position:The position to insert
  - value:The Bytes to insert
#### Returns:
  - error:Error if happened something
----



### Read(key:string) (Byte[], error or null):
#### Description
Read Bytes from the key
#### Args:
  - key:The Key to read
#### Returns:
  - value:The Bytes readed
  - error:Error if happened something
----
### ReadAt(key:string,position int64, size int64) (Byte[], error or null):
#### Description
Read Bytes from the key starting at position with the size
#### Args:
  - key:The Key to read
  - position:The position to start reading
  - size:The size of each chunck
#### Returns:
  - value:The Bytes readed
  - error:Error if happened something
----
### Delete (key: string) (error or null):
#### Description
Delete Bytes to the ke
#### Args:
  - key:The Key to insert
#### Returns:
  - error:Error if happened something
----
