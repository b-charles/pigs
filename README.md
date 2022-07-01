# Pigs

## What's Pigs?

Pigs is a Go library for bootstrapping micro-service applications.

## Why Pigs?

Because of [Pigs is Pigs](https://www.youtube.com/watch?v=GYXlF3sa9xs) very Goish conclusion:
> No more will I be a fool.
>
> Whenever comes to lifestock,
>
> Dash every single rule.
>
> If the animals come in singles
>
> Or if they come in sets,
>
> If they've got four feet and they're alive,
>
> They'll be classified as pet.
>
>
> -- <cite>Mike Flannery</cite>

## Where's Pigs?

You should be able to use Pigs with the `go get` command:
```
$ go get github.com/b-charles/pigs/<module>
```
with `<module>` the module you want to use.

## How's Pigs?

Pigs is collection of small libraires exposed it in an IOC framework.

The librairies are:
* ioc: The IOC framework. Completly written from scratch for the Pigs needs.
* config: Configuration collector.
* smartconf: Smart configuration caster.
* log: logging for the win.

## When Pigs?

TODO list:
* [Logrus](https://github.com/sirupsen/logrus) integration
* [Iris](https://github.com/kataras/iris) integration, with actuators:
    * components: complete list of components in IOC module (?)
    * env: complete configuration map
    * health
    * shutdown
* Unix signal support
* ApplicationContainer's life-cycle awareness
* [Cobra](https://github.com/spf13/cobra) integration (Maybe)

