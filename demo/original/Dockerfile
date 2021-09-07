FROM python:3.9-alpine

RUN pip install flask

COPY www data /

EXPOSE 5000

CMD [ "python", "/www/app.py" ]
