import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus } from '@nestjs/common';
import { ValidationError } from 'class-validator';

@Catch()
export class PrismaTransactionExceptionFilter implements ExceptionFilter {
    catch(exception: any, host: ArgumentsHost) {
        const ctx = host.switchToHttp();
        const response = ctx.getResponse();
        const request = ctx.getRequest();

        let status = HttpStatus.INTERNAL_SERVER_ERROR;
        let message: any = 'Internal server error';

        if (exception instanceof HttpException) {
            status = exception.getStatus();
            message = exception.getResponse();
        } else if (Array.isArray(exception.response?.message) && exception.response.message[0] instanceof ValidationError) {
            status = HttpStatus.UNPROCESSABLE_ENTITY;
            message = exception.response.message.map((err: ValidationError) => ({
                property: err.property,
                constraints: err.constraints,
            }));
        }

        response.status(status).json({
            message,
        });
    }
}
