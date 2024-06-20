import { ExceptionFilter, Catch, ArgumentsHost, UnprocessableEntityException } from '@nestjs/common';
import { Response } from 'express';

@Catch(UnprocessableEntityException)
export class UnprocessableEntityExceptionFilter implements ExceptionFilter {
    catch(exception: UnprocessableEntityException, host: ArgumentsHost) {
        const ctx = host.switchToHttp();
        const response = ctx.getResponse<Response>();
        const status = 422;

        response
            .status(status)
            .json({
                statusCode: status,
                message: exception.getResponse()['message'],
                error: 'Unprocessable Entity'
            });
    }
}
