# Sistema de generación de paquetes

Este proyecto es un generador de código con base a la arquitectura hexagonal, este proyecto gener el paquete,
los scripts de base de datos y los handlers referentes al modelo.

## Estructura

Para generar el paquete se tiene que seguir la siguiente estructura:

1. handler: Campo que indica si el paquete que se va a crear requerira crear los archivos API REST.
2. id: Indica el tipo de id para el paquete en la base de datos (`int`, `int64`, `uuid`).
3. name: Indica el nombre del paquete, este tiene que empezar por una letra del abecedario sin tildes y debe empezar con
   la primera letra en mayuscula y en singular.
4. db: Indica el nombre de la base de datos, esta compuesta por el nombre del schema y el nombre del paquete,
   ejemplo: `cfg.module`, donde `cfg` es el nombre del schema y `module` es el nombre del paquete (el nombre debe de
   estar en el formato snake_case y en plural).
5. fields: Lista los campos (atributos) que tendrá el paquete.

## Formato de los campos

#### ID

* **uuid**: uniqueidentifier,
* **int64**: –9223372036854775808 a 9223372036854775808,
* **int**: -2147483648 a 2147483648

#### Nombre del Paquete

* Digite el nombre del paquete en singular y CamelCase: `PaqueteName`.

#### Nombre de la Tabla

* Digite el nombre de la tabla en plural y snake_case con el nombre del schema: `dbo.table`

#### Campos

El formato es: `nombre:tipo:nonulo:tamaño`.

* **nombre**: Nombre del campo, minúsculas.
* **tipo**: `string`, `int`, `int64`, `float32`, `float64`, `time.Time`, `bool`, `uuid`.
* **nonulo**: `t` si permite nulos, `f` no permite nulos.
* **tamaño**: Número entero. Sólo aplica para `string` los demas deben ir con `0`.

Cada campo debe estar separada por un espacio, ejemplo:

  ````
  name:string:f:50 age:int:f:0 dni:int64:f:0 birth:time.Time:t:0 other:bool:t:0 role_id:uuid:f:0
  ````

## Configuración de variables de entorno

El servicio necesita la configuración de las variables de entorno para que pueda funcionar con total normalidad, toda
esta configuración se debe de realizar en el archivo con nombre único llamado ``config.json``:

````json
{
  "dest": "",
  "package_name": "",
  "src": ""
}
````

Donde:

* **dest**: Es la ruta donde se guardará el paquete.
* **package_name**: Es el nombre del paquete.
* **src**: Es la ruta donde se encuentra el archivo csv.

## Compilación y Ejecución

Para ejecutar el proyecto se tiene que seguir los siguientes pasos:
> _[NOTE]_ Para más información sobre cross compilation consultar la siguiente
> página [GO COMPILATION](https://golang.org/doc/install/source#environment)

1. Compilación

#### Sistema Operativo Linux

````shell
GOOS=linux GOARCH=amd64 go build -o go-wizard-open-jungle
````

#### Sistema Operativo Windows

````shell
GOOS=windows GOARCH=amd64 go build -o go-wizard-open-jungle
````

2. Ejecución

````shell
./go-wizard-open-jungle
````

## Ejemplo archivo csv

````csv
int;Module;cfg.modules;name:string:f:20 url_font:string:f:255 class:string:f:50 module_id:string:f:50
int64;Element;cfg.elements;name:string:f:20 url_font:string:f:255 class:string:f:50 module_id:string:f:50
uuid;Component;cfg.components;name:string:f:20 url_font:string:f:255 class:string:f:50 module_id:string:f:50
````
