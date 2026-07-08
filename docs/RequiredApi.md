

### Write (key: string, value: Byte[]) (error or null):
#### Description
Write Bytes to the key
#### Args:
  - key:The Key to insert
  - value:The Bytes to insert
#### Returns:
  - error:Error if happened something

----
### IncrementOrDecrement(key:string, value:int64) (int64, error or null):
#### Description
Increment or decrement the value of the key by value and return the new value. These function dont requires a locker, since it garantee atomicity 
#### Args:
  - key:The Key to increment
  - value:The value to increment or decrement (can be negative)
#### Returns:
  - value:The final value
  - error:Error if happened something
---- 
### WriteIfKeyNotExists(key: string, value: Byte[]) (error or null):
#### Description
Write Bytes to the key if the key not exists, this operations is atomic and doesn't require a locker
#### Args:
  - key:The Key to insert
  - value:The Bytes to insert
#### Returns:
  - error:Error if happened something

----
### WriteifValueEquals(key: string, value: Byte[], oldValue: Byte[]) (error or null):
#### Description
Write Bytes to the key if the value equals oldValue, This operation is atomic and doesn't require a locker
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
#### Exists(key: string) (bool, error or null):
#### Description
Check if the key exists
#### Args:
  - key:The Key to check
#### Returns:
  - exists:True if the key exists, false otherwise
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
### Lock(key:string, time int) (error or null):
#### Description
Lock the key for a specified time
#### Args:
  - key:The Key to lock
  - time:The time to live in seconds
#### Returns:
  - error:Error if happened something
----

### UnLock(key:string) (error or null):
#### Description
UnLock the key
#### Args:
  - key:The Key to unlock
#### Returns:
  - error:Error if happened something
----