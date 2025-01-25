# Multiple-Kind Search

Welcome to the **Multiple-Kind Search** project! This application serves as a versatile search engine designed to help users efficiently search through various types of data. Whether you're looking for news articles, author information, this project aims to provide a seamless search experience.

## Running the Application

To run the application, simply execute the following command in your terminal:

```bash
docker compose up -d
```

To seed the Elasticsearch database with initial data, run the following command:

```bash
sh /data/seed-elasticsearch.sh
```


## Project Aim

The primary goal of this project is to enable users to search for different kinds of data in a unified manner. By leveraging modern technologies, we aim to deliver fast and accurate search results, enhancing the overall user experience.

## Data Types Supported

This search engine supports the following data types:

- **Markdown + HTML based news content**: Search through rich content formatted in Markdown and HTML.
- **Authors of news**: Discover information about authors contributing to the news.

## Key Components

The project is built using a combination of powerful tools and technologies:

- **Docker Compose**: Simplifies the setup and management of the application environment.
- **Elasticsearch (version 7.10.2)**: A robust search engine that provides fast and scalable search capabilities, limited to this version for specific use cases.
- **Golang Backend**: Handles search queries and data processing efficiently.
- **React Frontend**: Offers a user-friendly interface for interacting with the search engine.

## Cursor Usage Details

This application was developed using Cursor IDE in just 5 hours. 
The rapid development process allowed for efficient coding and implementation of features, 
Showcasing the power of AI-assisted programming in accelerating project timelines.

## Note
__This project was designed by someone with no sense of aesthetics; I apologize for the ugliness of the designs.__