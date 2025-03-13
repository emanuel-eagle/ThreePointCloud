FROM public.ecr.aws/lambda/python:3.11

# Install the specified packages
RUN pip install boto3 beautifulsoup4 requests
# Copy function code
COPY src/careerStats/careerStats.py ${LAMBDA_TASK_ROOT}

# Set the CMD to your handler
CMD [ "careerStats.handler" ]