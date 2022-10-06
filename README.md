# mtg

## Preparation

```
$ docker-compose -f postgres.yaml up
```

## How to use

You will get empty table when you start the program.

To import cards into table use http://Localhost:9001/import/NUM
where NUM is page number from which you want to import the cards.
It will import all cards from selected page.

For listing use  http://Localhost:9001/list?page=NUM
It will return one page from all cards in a table.

For Searching use  http://Localhost:9001/list?PARAMS 
PARAMS should always have page number and it should be on last spot.

Some examples:
* "colors=R&page=3"
* "colors=R&rarity=Rare&page=2"              

For searching colors use B for black, R for red, U for blue, W for white, G for green

If you have ID of the card you can use it to get all information on that card.
Use http://Localhost:9001/card/ID.

