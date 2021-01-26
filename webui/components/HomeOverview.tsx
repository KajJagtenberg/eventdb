import { Flex, Heading, Text } from '@chakra-ui/react';

const HomeOverview = ({ version, size, human_size }) => {
  return (
    <Flex
      direction="column"
      p={4}
      m={4}
      bg="white"
      style={{
        boxShadow: '0 0 10px #aaa',
      }}
    >
      <Heading size="md" color="gray.600" textTransform="uppercase" my={2}>
        Info
      </Heading>

      <Flex>
        <Text fontWeight="semibold" mr={4}>
          Version:
        </Text>
        <Text>{version}</Text>
      </Flex>

      <Flex>
        <Text fontWeight="semibold" mr={4}>
          Size human:
        </Text>
        <Text>{human_size}</Text>
      </Flex>

      <Flex>
        <Text fontWeight="semibold" mr={4}>
          Size bytes:
        </Text>
        <Text>{size}</Text>
      </Flex>
    </Flex>
  );
};

export default HomeOverview;
