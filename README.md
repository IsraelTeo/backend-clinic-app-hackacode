***API - Clínica Médica***

Una API RESTful orientada a la gestión de una bodega.

**Este proyecto contiene:**

    Esta aplicación sirve como backend para una clínica médica y fue desarrollada como parte del hackathon "Hackacode" organizado por 
    Todocode. Fue presentada el 25 de febrero de 2025 y, desde entonces, ha sido ampliada con nuevas funcionalidades, importantes mejoras, 
    refactorizaciones y actualizaciones. El sistema está construido utilizando el lenguaje de programación Go y permite realizar operaciones 
    CRUD para gestionar pacientes, doctores y citas médicas, asegurando una correcta programación basada en la disponibilidad de los doctores. 
    Además, cuenta con un sistema de inicio de sesión basado en autorización.

**Tecnologías Usadas**

    Java 17
     Lenguaje de programación orientado a objetos, robusto y ampliamente utilizado en el desarrollo empresarial por su seguridad.

    Spring Boot
    Framework para crear aplicaciones backend de forma rápida y con configuración mínima.

    Spring Data JPA
    Módulo de Spring Framework que funciona como una abstracción sobre JPA para facilitar la persistencia 
    de datos.

    PostgreSQL
    Base de datos relacional ideal para aplicaciones empresariales.

    Lombok
    Elimina código repetitivo como getters, setters, builders, etc., mediante anotaciones.

    MapStruct
    Framework para el mapeo automático entre entidades y DTOs, basado en interfaces.

    Spring Validation
    Validación declarativa de campos en DTOs usando anotaciones.

    Swagger / OpenAPI
    Documentación interactiva de endpoints REST directamente desde el código.

    SLF4J + Logback
    Logging flexible y personalizable para auditoría y depuración.

 **Buenas prácticas**

    Uso de archivos .yml para centralizar configuraciones como puertos, credenciales de base de datos, etc.

    Inyección de dependencias a través de constructores utilizando Lombok (@RequiredArgsConstructor) 
    para evitar acoplamiento directo con el framework y facilitar las futuras pruebas unitarias.

    Manejo centralizado de excepciones.

    Documentación automática y actualizada de los endpoints REST con Swagger / OpenAPI.

    Uso de DTOs para transferir datos entre cliente y servidor.

    Validación de datos de entrada en los DTOs utilizando anotaciones como @NotBlank, @Size, 
    @Valid, entre otras.

    Implementación de paginación en endpoints que devuelven listas límitadas de datos para 
    mejorar rendimiento y escalabilidad.

    Arquitectura de capas bien definida: controller, service, repository, dto, mapper, 
    exception, domain, mapper, criteria.

    Uso de nombres descriptivos para variables, constantes, métodos, clases, paquetes 
    e interfaces.

    Aplicación del principio de responsabilidad única (SRP) del conjunto SOLID para mantener 
    clases y métodos mantenibles.

    Mapeos limpios y desacoplados entre entidades y DTOs utilizando MapStruct.

    Logging estructurado con SLF4J y Logback para monitorear el flujo y los errores.

  **Certificado de Participación:**
![hackaton-1](https://github.com/user-attachments/assets/5e8854ab-4302-4763-a4c2-816a2575d85b)

 **Presentación oficial – Hackacode (25 de febrero de 2025):**

**Puedes ver la presentación oficial de los proyectos participantes en el siguiente enlace:**
🔗 Ver en YouTube

Israel Juárez (yo) – Grupo 1
🕒 Hora: 3h     📍 Minuto: 31
