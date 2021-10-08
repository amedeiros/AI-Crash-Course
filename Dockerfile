FROM python:3.7

WORKDIR /usr/src/app

COPY requirements.txt ./
# COPY requirements-dev.txt ./
RUN apt update
RUN apt upgrade -y
RUN apt install build-essential -y
RUN python -m pip install --upgrade pip

# Install the rest of the requirements
RUN python -m pip install -r requirements.txt

COPY . .
