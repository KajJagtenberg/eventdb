import {
  Flex,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Link as UILink,
} from '@chakra-ui/react';
import { gql, useQuery } from '@apollo/client';
import Link from 'next/link';
import dayjs from 'dayjs';

export const StreamTable = ({ id }) => {
  const { data } = useQuery(
    gql`
      query Get($stream: String!) {
        get(input: { stream: $stream, version: 0, limit: 0 }) {
          id
          version
          type
          data
          metadata
          causation_id
          correlation_id
          added_at
        }
      }
    `,
    {
      variables: {
        stream: id,
      },
    }
  );

  return (
    <Flex>
      <Table variant="striped" size="sm">
        <Thead>
          <Tr>
            <Th>ID</Th>
            <Th>Version</Th>
            <Th>Type</Th>
            <Th>Data</Th>
            <Th>Metadata</Th>
            <Th>Causation ID</Th>
            <Th>Correlation ID</Th>
            <Th>Added At</Th>
          </Tr>
        </Thead>
        <Tbody>
          {data?.get.map(
            (
              {
                id,
                version,
                type,
                data,
                metadata,
                causation_id,
                correlation_id,
                added_at,
              }: {
                id: string;
                version: number;
                type: string;
                data: string;
                metadata: string;
                causation_id: string;
                correlation_id: string;
                added_at: number;
              },
              index: number
            ) => (
              <Tr key={index}>
                <Td>
                  <UILink color="blue.400">
                    <Link href={`/events/${id}`}>
                      <a>{id}</a>
                    </Link>
                  </UILink>
                </Td>
                <Td>{version}</Td>
                <Td>{type}</Td>
                <Td>{atob(data)}</Td>
                <Td>{atob(metadata) || '-'}</Td>
                <Td>
                  <UILink color="blue.400">
                    <Link href={`/events/${causation_id}`}>
                      <a>{causation_id}</a>
                    </Link>
                  </UILink>
                </Td>
                <Td>
                  <UILink color="blue.400">
                    <Link href={`/events/${correlation_id}`}>
                      <a>{correlation_id}</a>
                    </Link>
                  </UILink>
                </Td>
                <Td>{dayjs(added_at).format('HH:mm:ss DD-MM-YYYY')}</Td>
              </Tr>
            )
          )}
        </Tbody>
      </Table>
    </Flex>
  );
};
