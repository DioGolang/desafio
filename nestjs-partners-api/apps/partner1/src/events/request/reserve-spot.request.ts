import { TicketKind } from '@prisma/client';
import {IsArray, IsEmail, isEmail, IsIn, IsNotEmpty, IsString,} from "class-validator";

export class ReserveSpotRequest {
  @IsNotEmpty({message: 'the spots field is empty'})
  @IsString({each: true})
  spots: string[]; //['A1', 'A2']

  @IsNotEmpty({message: 'the spots field is empty'})
  @IsIn(['full', 'half'], { message: 'The type must be either "full" or "half"' })
  ticket_kind: TicketKind;

  @IsEmail()
  email: string;
}
