import { Flex, Heading, Spacer, Text } from '@chakra-ui/react';
import Link from 'next/link';

const links = [
  {
    href: '/dashboard',
    text: 'Dashboard',
  },
  {
    href: '/streams',
    text: 'Streams',
  },
  {
    href: '/projections',
    text: 'Projections',
  },
];

export const Navbar = () => {
  return (
    <Flex bg="brand.600" h={10} color="white">
      <Flex alignItems="center" mx={4}>
        <Heading size="md">
          <Link href="/dashboard">
            <a>EventflowDB</a>
          </Link>
        </Heading>
      </Flex>

      <Spacer />

      <Flex alignItems="center">
        {links.map(({ text, href }, index) => (
          <Link href={href} key={index}>
            <a>
              <Text fontWeight="semibold" mx={2} fontSize={14}>
                {text}
              </Text>
            </a>
          </Link>
        ))}
      </Flex>
    </Flex>
  );
};
