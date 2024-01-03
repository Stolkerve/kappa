# Kappa

Kappa o κ, es la letra que antecede λ en el alfabeto griego. Por que kappa es 1000% mejor que aws lambda 😎.
Kappa basicamente es una servicio serverless basico basado en aws lambda que ejecuta funciones compiladas en WASM. Gracias a WASM, puedes crear funciones en distintos lenguajes de programacion con un binario cercano a codigo de maquina.

# Como usar Kappa
Es muy facil de usar, puedes interactuar con el servidor con su api
| ruta | metodo | descripción |
|---|---|---|
| /api/deploy | POST | Espera un formulacion con el archivo WASM con el nombre de function y un tamaño maximo de 5MB. Retorna el id de funcion. Texto plano |
| /api/deploys/{page:[0-9]+} | GET | Retorna una lista de ids de las funciones que han sido subidas con un maximo de 50 elementos. Formato JSON |
| /api/call/{id}/* | ANY | Ejecuta la funcion... No importa el metodo. |
| /api/calls/{id}/{page:[0-9]+} | GET | Retorna una lista de las ejecuciones de las funciones con un maximo de 50 elementos. Formato JSON|