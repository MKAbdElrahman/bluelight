## Tips

- Don't create a new package untill you suffer missing the organization it adds.
- Put dependancies closer to where its needed.



## Rememeber

- The domain entity have the logic to validate itself always being in a a valid state
    - it errors at creation or mutation time if needed

- I tried adding validation logic to the contracts layer but it got messy

- contracts are responsible for ensuring data is coming and going in the correct type
    - they return bad requests

- entities enforce business rules like "title can't be empty"
    - they return unprocessible entites


## To-Do List

[x] group movie handlers under a parent package.

[x] refactor the error hadning and validation for the movies resource
[x] return one validation error at a time 
[x] move validation logic to the new movie factory function

[x] list movies feature: here we should move the query validation to the application service object  from the contract.  