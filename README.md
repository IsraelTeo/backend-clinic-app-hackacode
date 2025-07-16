***API - Clínica Médica***

Una API RESTful orientada a la gestión de una bodega.

**Este proyecto contiene:**

    Esta aplicación sirve como backend para una clínica médica y fue desarrollada como parte del hackathon 
    "Hackacode" organizado por Todocode. Fue presentada el 25 de febrero de 2025 y, desde entonces, 
    ha sido  ampliada con nuevas funcionalidades, importantes mejoras, refactorizaciones y actualizaciones. 
    
    El sistema está construido utilizando el lenguaje de programación Go y permite realizar operaciones 
    CRUD para gestionar pacientes, doctores y citas médicas, asegurando una correcta programación basada 
    en la disponibilidad de los doctores. 
    
    Además, cuenta con un sistema de inicio de sesión basado en autorización.

Tecnologías y herramientas

    Go: Lenguaje de programación principal utilizado para el desarrollo del backend.

    MySQL: Sistema de gestión de bases de datos relacional.

    GORM: ORM para interactuar con bases de datos relacionales en Go.

    JWT: Utilizado para autenticación de usuarios y gestión segura del acceso.

    Echo Framework: Framework web de alto rendimiento para construir APIs RESTful en Go.

    Docker: Herramienta de contenedorización que facilita el despliegue y la consistencia 
    del entorno.
    
Librerías:
    
    github.com/go-playground/locales: Proporciona soporte de localización (idiomas, formatos de fecha, etc.).
    
    github.com/go-playground/universal-translator: Motor de traducción usado para internacionalización y 
    validaciones de mensajes.
    
    github.com/go-playground/validator / validator/v10: Validación estructural de structs y campos 
    (como @Valid en Java).
    
    github.com/joho/godotenv: Carga variables de entorno desde archivos .env. Muy útil en desarrollo.
    
    github.com/jung-kurt/gofpdf: Generación de archivos PDF directamente desde Go.
    
    golang.org/x/crypto: Conjunto de algoritmos y utilidades criptográficas recomendadas por el equipo de Go.
    

 **Buenas prácticas**

    Inyección de dependencias a través de constructores.

    Manejo explícito de errores.

    Validación de datos de entrada utilizando la librería: github.com/go-playground/validator

    Implementación de paginación en endpoints que devuelven listas límitadas de datos para 
    mejorar rendimiento y escalabilidad.

    Arquitectura de capas bien definida: handler, logic(service), repository, model.

    Uso de nombres descriptivos para variables, constantes, métodos, structs, paquetes 
    e interfaces.

    Aplicación del principio de responsabilidad única (SRP) del conjunto SOLID para mantener 
    paquetes y métodos mantenibles.
    
  **Certificado de Participación:**
![hackaton-1](https://github.com/user-attachments/assets/5e8854ab-4302-4763-a4c2-816a2575d85b)

 **Presentación oficial – Hackacode (25 de febrero de 2025):**

**Puedes ver la presentación oficial de los proyectos participantes en el siguiente enlace:**
🔗 [Ver en YouTube](https://www.youtube.com/watch?v=Nr6f0MuI_rM&t=13426s)

Israel Juárez (yo) – Grupo 1
🕒 Hora: 3h     📍 Minuto: 31
