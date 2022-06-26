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

You should be able to install Pigs with using dep:
```
$ go get github.com/b-charles/pigs
```
or something. I don't know, I've never tried yet.

## How's Pigs?

Pigs is collection of small libraires, for the most only wrapping an existing well done librairy and expose it in an IOC framework.

The librairies are:
* ioc: The IOC framework. Completly written from scratch for the Pigs needs.
* filesystem: Wrapping of [afero](https://github.com/spf13/afero)
* config: Configuration collector.

## When Pigs?

One day, maybe. But not before the end of this TODO list:
* [Logrus](https://github.com/sirupsen/logrus) integration
* [Iris](https://github.com/kataras/iris) integration, with actuators:
    * components: complete list of components in IOC module (?)
    * env: complete configuration map
    * health
    * loggers: show and modifies the condfiguation of logger (?)
    * shutdown
* Unix signal support
* [Cobra](https://github.com/spf13/cobra) integration (Maybe)
* [Viper](https://github.com/spf13/viper) integration (As repacement of configuration? Second level of configuration with dynamic features?)

