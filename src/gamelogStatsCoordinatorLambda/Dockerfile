FROM public.ecr.aws/lambda/python:3.11

# Install the specified packages
RUN pip install boto3
# Copy function code
COPY src/gamelogStatsCoordinatorLambda/coordinator.py ${LAMBDA_TASK_ROOT}

# Set the CMD to your handler
CMD [ "coordinator.handler" ]